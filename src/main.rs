#[macro_use]
extern crate clap;
extern crate exitcode;
extern crate term;

mod cli;
mod config;
mod terminal;

fn main() {
    let app = cli::CLI::new();
    std::process::exit(app.run());
}
