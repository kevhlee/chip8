use rand::{rngs::ThreadRng, Rng};

use super::{display::Display, keyboard::Keyboard, sound::Sound, timer::Timer};

fn decode_x(opcode: u16) -> usize {
    ((opcode >> 8) & 0xF) as usize
}

fn decode_y(opcode: u16) -> usize {
    ((opcode >> 4) & 0xF) as usize
}

fn decode_addr(opcode: u16) -> u16 {
    opcode & 0xFFF
}

fn decode_byte(opcode: u16) -> u8 {
    (opcode & 0xFF) as u8
}

fn decode_nibb(opcode: u16) -> usize {
    (opcode & 0xF) as usize
}

const FONTS: [u8; 0x50] = [
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80, // F
];

const FONT_SIZE: u16 = 5;
const PROGRAM_START_ADDRESS: usize = 0x200;

pub struct VirtualMachine {
    i: u16,
    pc: u16,
    sp: u8,
    v: [u8; 0x10],
    memory: [u8; 0x1000],
    stack: [u16; 0x10],
    rng: ThreadRng,
}

impl VirtualMachine {
    pub fn new() -> Self {
        let mut vm = Self {
            i: 0,
            pc: PROGRAM_START_ADDRESS as u16,
            sp: 0,
            v: [0; 0x10],
            memory: [0; 0x1000],
            stack: [0; 0x10],
            rng: rand::rng(),
        };

        vm.reset();
        vm
    }

    pub fn reset(&mut self) {
        self.i = 0;
        self.sp = 0;
        self.pc = PROGRAM_START_ADDRESS as u16;
        self.memory[..0x50].copy_from_slice(&FONTS);
    }

    pub fn load_bytes(&mut self, bytes: &Vec<u8>) -> Result<(), String> {
        if self.memory.len() < bytes.len() {
            return Err(String::from("Too many bytes"));
        }

        let mut index: usize = PROGRAM_START_ADDRESS;
        for byte in bytes {
            self.memory[index] = *byte;
            index += 1;
        }
        Ok(())
    }

    pub fn step(
        &mut self,
        timer: &mut Timer,
        sound: &mut Sound,
        display: &mut Display,
        keyboard: &mut Keyboard,
    ) -> Result<(), String> {
        let opcode = self.fetch_opcode();

        match opcode >> 12 {
            0x0 => self.exec_opcode_0x0(opcode, display),
            0x1 => self.exec_opcode_0x1(opcode),
            0x2 => self.exec_opcode_0x2(opcode),
            0x3 => self.exec_opcode_0x3(opcode),
            0x4 => self.exec_opcode_0x4(opcode),
            0x5 => self.exec_opcode_0x5(opcode),
            0x6 => self.exec_opcode_0x6(opcode),
            0x7 => self.exec_opcode_0x7(opcode),
            0x8 => self.exec_opcode_0x8(opcode),
            0x9 => self.exec_opcode_0x9(opcode),
            0xA => self.exec_opcode_0xA(opcode),
            0xB => self.exec_opcode_0xB(opcode),
            0xC => self.exec_opcode_0xC(opcode),
            0xD => self.exec_opcode_0xD(opcode, display),
            0xE => self.exec_opcode_0xE(opcode, keyboard),
            0xF => self.exec_opcode_0xF(opcode, timer, sound, keyboard),
            _ => unreachable!(),
        }
    }

    fn fetch_opcode(&mut self) -> u16 {
        let byte_hi = self.memory[self.pc as usize] as u16;
        let byte_lo = self.memory[(self.pc + 1) as usize] as u16;
        self.pc += 2;
        (byte_hi << 8) | byte_lo
    }

    fn exec_opcode_0x0(&mut self, opcode: u16, display: &mut Display) -> Result<(), String> {
        match opcode {
            // 00E0 - CLS
            0x00E0 => display.clear(),

            // 00EE - RET
            0x00EE => {
                if self.sp == 0 {
                    return Err(String::from("Empty call stack"));
                }
                self.sp -= 1;
                self.pc = self.stack[self.sp as usize];
            }

            _ => {}
        }

        Ok(())
    }

    fn exec_opcode_0x1(&mut self, opcode: u16) -> Result<(), String> {
        // 1nnn - JP addr
        self.pc = decode_addr(opcode);
        Ok(())
    }

    fn exec_opcode_0x2(&mut self, opcode: u16) -> Result<(), String> {
        // 2nnn - CALL addr
        if self.sp > 0xF {
            Err(String::from("Call stack overflow"))
        } else {
            self.stack[self.sp as usize] = self.pc;
            self.sp += 1;
            self.pc = decode_addr(opcode);
            Ok(())
        }
    }

    fn exec_opcode_0x3(&mut self, opcode: u16) -> Result<(), String> {
        // 3xkk - SE Vx, byte
        if self.v[decode_x(opcode)] == decode_byte(opcode) {
            self.pc += 2;
        }
        Ok(())
    }

    fn exec_opcode_0x4(&mut self, opcode: u16) -> Result<(), String> {
        // 4xkk - SNE Vx, byte
        if self.v[decode_x(opcode)] != decode_byte(opcode) {
            self.pc += 2;
        }
        Ok(())
    }

    fn exec_opcode_0x5(&mut self, opcode: u16) -> Result<(), String> {
        // 5xy0 - SE Vx, Vy
        if decode_nibb(opcode) == 0 && self.v[decode_x(opcode)] == self.v[decode_y(opcode)] {
            self.pc += 2;
        }
        Ok(())
    }

    fn exec_opcode_0x6(&mut self, opcode: u16) -> Result<(), String> {
        // 6xkk - LD Vx, byte
        self.v[decode_x(opcode)] = decode_byte(opcode);
        Ok(())
    }

    fn exec_opcode_0x7(&mut self, opcode: u16) -> Result<(), String> {
        // 7xkk - LD Vx, byte
        self.v[decode_x(opcode)] = self.v[decode_x(opcode)].wrapping_add(decode_byte(opcode));
        Ok(())
    }

    fn exec_opcode_0x8(&mut self, opcode: u16) -> Result<(), String> {
        let x = decode_x(opcode);
        let y = decode_y(opcode);

        match decode_nibb(opcode) {
            // 8xy0 - LD Vx, Vy
            0x0 => self.v[x] = self.v[y],

            // 8xy1 - OR Vx, Vy
            0x1 => self.v[x] |= self.v[y],

            // 8xy2 - AND Vx, Vy
            0x2 => self.v[x] &= self.v[y],

            // 8xy3 - XOR Vx, Vy
            0x3 => self.v[x] ^= self.v[y],

            // 8xy4 - ADD Vx, Vy
            0x4 => {
                let (result, overflow) = self.v[x].overflowing_add(self.v[y]);
                self.v[0xF] = if overflow { 1 } else { 0 };
                self.v[x] = result;
            }

            // 8xy5 - SUB Vx, Vy
            0x5 => {
                let (result, borrow) = self.v[x].overflowing_sub(self.v[y]);
                self.v[0xF] = if borrow { 0 } else { 1 };
                self.v[x] = result;
            }

            // 8xy6 - SHR Vx {, Vy}
            0x6 => {
                self.v[0xF] = self.v[x] & 1;
                self.v[x] >>= 1;
            }

            // 8xy7 - SUBN Vx, Vy
            0x7 => {
                let (result, borrow) = self.v[y].overflowing_sub(self.v[x]);
                self.v[0xF] = if borrow { 0 } else { 1 };
                self.v[x] = result;
            }

            // 8xyE - SHL Vx {, Vy}
            0xE => {
                self.v[0xF] = self.v[x] >> 7;
                self.v[x] <<= 1;
            }

            _ => {}
        }

        Ok(())
    }

    fn exec_opcode_0x9(&mut self, opcode: u16) -> Result<(), String> {
        // 9xy0 - SNE Vx, Vy
        if decode_nibb(opcode) == 0 && self.v[decode_x(opcode)] != self.v[decode_y(opcode)] {
            self.pc += 2;
        }
        Ok(())
    }

    #[allow(non_snake_case)]
    fn exec_opcode_0xA(&mut self, opcode: u16) -> Result<(), String> {
        // Annn - LD I, addr
        self.i = decode_addr(opcode);
        Ok(())
    }

    #[allow(non_snake_case)]
    fn exec_opcode_0xB(&mut self, opcode: u16) -> Result<(), String> {
        // Bnnn - JP V0, addr
        self.pc = decode_addr(opcode).wrapping_add(self.v[0] as u16);
        Ok(())
    }

    #[allow(non_snake_case)]
    fn exec_opcode_0xC(&mut self, opcode: u16) -> Result<(), String> {
        // Cxkk - RND Vx, byte
        self.v[decode_x(opcode)] = self.rng.random::<u8>() & decode_byte(opcode);
        Ok(())
    }

    #[allow(non_snake_case)]
    fn exec_opcode_0xD(&mut self, opcode: u16, display: &mut Display) -> Result<(), String> {
        // Dxyn - DRW Vx, Vy, nibb
        let x = self.v[decode_x(opcode)] as usize;
        let y = self.v[decode_y(opcode)] as usize;

        let mut flag = false;
        for n in 0..decode_nibb(opcode) {
            let mut byte = self.memory[(self.i as usize) + n];

            for i in (0..8).rev() {
                if display.set(x + i, y + n, (byte & 0x01) == 1) {
                    flag = true;
                }
                byte >>= 1;
            }
        }

        self.v[0xF] = if flag { 1 } else { 0 };

        Ok(())
    }

    #[allow(non_snake_case)]
    fn exec_opcode_0xE(&mut self, opcode: u16, keyboard: &mut Keyboard) -> Result<(), String> {
        let x = decode_x(opcode);

        match decode_byte(opcode) {
            // Ex9E - SKP Vx
            0x9E => {
                if keyboard.is_pressed(self.v[x]) {
                    self.pc += 2;
                }
            }

            // ExA1 - SKNP Vx
            0xA1 => {
                if !keyboard.is_pressed(self.v[x]) {
                    self.pc += 2;
                }
            }

            _ => {}
        }

        Ok(())
    }

    #[allow(non_snake_case)]
    fn exec_opcode_0xF(
        &mut self,
        opcode: u16,
        timer: &mut Timer,
        sound: &mut Sound,
        keyboard: &mut Keyboard,
    ) -> Result<(), String> {
        let x = decode_x(opcode);

        match decode_byte(opcode) {
            // Fx07 - LD Vx, DT
            0x07 => self.v[x] = timer.read(),

            // Fx0A - LD Vx, K
            0x0A => {
                let (key, polled) = keyboard.poll();
                if polled {
                    self.v[x] = key;
                } else {
                    self.pc -= 2;
                }
            }

            // Fx15 - LD DT, Vx
            0x15 => timer.write(self.v[x]),

            // Fx18 - LD ST, Vx
            0x18 => sound.write(self.v[x]),

            // Fx1E - ADD I, Vx
            0x1E => self.i += self.v[x] as u16,

            // Fx29 - LD F, Vx
            0x29 => self.i = (self.v[x] as u16) * FONT_SIZE,

            // Fx33 - LD B, Vx
            0x33 => {
                self.memory[self.i as usize] = self.v[x] / 100;
                self.memory[(self.i + 1) as usize] = (self.v[x] % 100) / 100;
                self.memory[(self.i + 2) as usize] = self.v[x] % 10;
            }

            // Fx55 - LD [I], Vx
            0x55 => {
                for i in 0..=x {
                    self.memory[self.i as usize + i] = self.v[i];
                }
            }

            // Fx65 - LD Vx, [I]
            0x65 => {
                for i in 0..=x {
                    self.v[i] = self.memory[self.i as usize + i];
                }
            }

            _ => {}
        }

        Ok(())
    }

    #[allow(dead_code)]
    fn print_debug_info(&self, opcode: u16) {
        println!("==== Opcode: {:04X} ====", opcode);

        println!(" I: {:03X}", self.i);
        println!("PC: {:03X}", self.pc);

        println!();
        for i in 0..(self.sp as usize) {
            print!("[ {:03X} ]", self.stack[i]);
        }
        println!();
        println!();

        for i in 0..self.v.len() {
            print!(" {:X}: {:02X} ", i, self.v[i]);
            if i % 4 == 3 {
                println!();
            }
        }
    }
}
