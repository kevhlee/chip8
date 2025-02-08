pub struct Timer {
    value: u8,
}

impl Timer {
    pub fn new() -> Self {
        Self { value: 0 }
    }

    pub fn step(&mut self) {
        if self.value > 0 {
            self.value -= 1;
        }
    }

    pub fn read(&self) -> u8 {
        self.value
    }

    pub fn write(&mut self, value: u8) {
        self.value = value;
    }
}
