pub struct Sound {
    value: u8,
}

impl Sound {
    pub fn new() -> Self {
        Self { value: 0 }
    }

    pub fn write(&mut self, value: u8) {
        self.value = value;
    }

    pub fn step(&mut self) {
        if self.value > 0 {
            self.value -= 1;
        }
    }
}
