use sdl2::{pixels::Color, rect::Rect, render::WindowCanvas};

pub const DISPLAY_WIDTH: usize = 0x40;
pub const DISPLAY_HEIGHT: usize = 0x20;

const FOREGROUND: Color = Color::RGB(0xFF, 0xFF, 0xFF);
const BACKGROUND: Color = Color::RGB(0x00, 0x00, 0x00);

pub struct Display {
    rect: Rect,
    scale: u32,
    canvas: WindowCanvas,
    buffer: [bool; DISPLAY_WIDTH * DISPLAY_HEIGHT],
}

impl Display {
    pub fn new(canvas: WindowCanvas, scale: u32) -> Self {
        Self {
            rect: Rect::new(0, 0, scale, scale),
            scale,
            canvas,
            buffer: [false; DISPLAY_WIDTH * DISPLAY_HEIGHT],
        }
    }

    pub fn clear(&mut self) {
        self.buffer.fill(false);
    }

    pub fn set(&mut self, x: usize, y: usize, value: bool) -> bool {
        let idx = ((y % DISPLAY_HEIGHT) * DISPLAY_WIDTH) + (x % DISPLAY_WIDTH);
        let flag = self.buffer[idx] & value;
        self.buffer[idx] ^= value;
        return flag;
    }

    pub fn render(&mut self) {
        self.canvas.set_draw_color(BACKGROUND);
        self.canvas.clear();
        self.canvas.set_draw_color(FOREGROUND);

        for i in 0..self.buffer.len() {
            if !self.buffer[i] {
                continue;
            }

            self.rect.x = ((i % DISPLAY_WIDTH) as i32) * (self.scale as i32);
            self.rect.y = ((i / DISPLAY_WIDTH) as i32) * (self.scale as i32);

            let _ = self.canvas.fill_rect(self.rect);
        }

        self.canvas.present();
    }
}
