mod display;
mod keyboard;
mod sound;
mod timer;
mod vm;

use std::collections::HashMap;
use std::time::Duration;

use display::{Display, DISPLAY_HEIGHT, DISPLAY_WIDTH};
use keyboard::Keyboard;
use sdl2::{event::Event, keyboard::Keycode, pixels::Color};
use sound::Sound;
use timer::Timer;
use vm::VirtualMachine;

pub struct Options {
    pub scale: u32,
    pub tps: usize,
}

impl Options {
    pub fn defaults() -> Self {
        Self { scale: 8, tps: 12 }
    }
}

pub fn start(filename: String, opts: &Options) -> Result<(), String> {
    let context = sdl2::init()?;
    let video = context.video()?;

    let window = video
        .window(
            "CHIP-8",
            (DISPLAY_WIDTH as u32) * opts.scale,
            (DISPLAY_HEIGHT as u32) * opts.scale,
        )
        .position_centered()
        .opengl()
        .build()
        .map_err(|err| err.to_string())?;

    let mut canvas = window
        .into_canvas()
        .software()
        .build()
        .map_err(|err| err.to_string())?;

    canvas.set_draw_color(Color::RGB(0x00, 0x00, 0x00));
    canvas.clear();
    canvas.present();

    let mut event_pump = context.event_pump()?;

    let mut vm = VirtualMachine::new();
    let mut timer = Timer::new();
    let mut sound = Sound::new();
    let mut display = Display::new(canvas, opts.scale);
    let mut keyboard = Keyboard::new();

    let keymap: HashMap<Keycode, u8> = HashMap::from([
        (Keycode::Num1, 0x0),
        (Keycode::Num2, 0x1),
        (Keycode::Num3, 0x2),
        (Keycode::Num4, 0x3),
        (Keycode::Q, 0x4),
        (Keycode::W, 0x5),
        (Keycode::E, 0x6),
        (Keycode::R, 0x7),
        (Keycode::A, 0x8),
        (Keycode::S, 0x9),
        (Keycode::D, 0xA),
        (Keycode::F, 0xB),
        (Keycode::Z, 0xC),
        (Keycode::X, 0xD),
        (Keycode::C, 0xE),
        (Keycode::V, 0xF),
    ]);

    match std::fs::read(filename) {
        Ok(bytes) => {
            if let Err(error) = vm.load_bytes(&bytes) {
                return Err(error);
            }
        }
        Err(error) => {
            return Err(error.to_string());
        }
    }

    'running: loop {
        for event in event_pump.poll_iter() {
            match event {
                Event::Quit { .. } => break 'running,

                Event::KeyDown {
                    keycode: Some(keycode),
                    ..
                } => match keycode {
                    Keycode::Escape => break 'running,

                    _ => {
                        if let Some(key) = keymap.get(&keycode) {
                            keyboard.set(*key, true);
                        }
                    }
                },

                Event::KeyUp {
                    keycode: Some(keycode),
                    ..
                } => {
                    if let Some(key) = keymap.get(&keycode) {
                        keyboard.set(*key, false);
                    }
                }

                _ => {}
            }
        }

        timer.step();
        sound.step();

        for _ in 0..opts.tps {
            match vm.step(&mut timer, &mut sound, &mut display, &mut keyboard) {
                Ok(_) => {}
                Err(error) => {
                    println!("Error: {}", error);
                }
            }
        }

        display.render();

        std::thread::sleep(Duration::new(0, 1_000_000_000u32 / 60));
    }

    Ok(())
}
