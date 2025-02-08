mod chip8;

use chip8::Options;

fn main() -> Result<(), String> {
    let path = std::env::args().nth(1).expect("no ROM path given");

    chip8::start(path, &Options::defaults())
}
