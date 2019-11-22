pub struct Terminal {
    stdout: Box<term::StdoutTerminal>,
    stderr: Box<term::StderrTerminal>,
}

impl Terminal {
    pub fn new() -> Self {
        Self {
            stdout: term::stdout().unwrap(),
            stderr: term::stderr().unwrap(),
        }
    }

    pub fn new_line(&mut self) {
        writeln!(self.stdout, "").unwrap();
    }

    pub fn error(&mut self, text: &str) {
        self.stderr.fg(term::color::RED).unwrap();
        writeln!(self.stderr, "{}", text).unwrap();
        self.stderr.reset().unwrap();
    }

    pub fn write(&mut self, text: &str, color: Option<term::color::Color>) {
        self.stdout.fg(color.unwrap_or(term::color::WHITE)).unwrap();
        write!(self.stdout, "{}", text).unwrap();
        self.stdout.reset().unwrap();
    }

    pub fn write_line(&mut self, text: &str) {
        self.write(text, None);
        self.new_line();
    }

    pub fn title(&mut self, text: &str) {
        self.write(text, Some(term::color::CYAN));
        self.new_line();
    }
}
