pub struct Keyboard {
    keys: [bool; 0x10],
    polling: bool,
    last_pressed: i16,
}

impl Keyboard {
    pub fn new() -> Self {
        Self {
            keys: [false; 0x10],
            polling: false,
            last_pressed: -1,
        }
    }

    pub fn set(&mut self, key: u8, value: bool) {
        if value && self.last_pressed < 0 {
            self.last_pressed = key as i16
        }
        self.keys[key as usize] = value;
    }

    pub fn is_pressed(&self, key: u8) -> bool {
        self.keys[key as usize]
    }

    pub fn poll(&mut self) -> (u8, bool) {
        if self.polling && self.last_pressed >= 0 {
            let key = self.last_pressed as u8;
            self.polling = false;
            self.last_pressed = -1;
            (key, true)
        } else {
            self.polling = true;
            (0, false)
        }
    }
}
