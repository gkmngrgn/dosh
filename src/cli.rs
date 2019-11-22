use crate::config::{Config, ConfigError, ENV_NAME, FILE_NAME};
use crate::terminal::Terminal;
use clap::{App, Arg, ArgMatches};
use std::env;
use std::fs::File;
use std::io::prelude::*;
use std::path::Path;

pub struct CLI<'a> {
    arg_matches: ArgMatches<'a>,
}

impl<'a> CLI<'a> {
    pub fn new() -> CLI<'a> {
        let arg_matches = App::new("Just do it, shh!")
            .version(crate_version!())
            .author(crate_authors!())
            .arg(
                Arg::with_name("init")
                    .short("i")
                    .long("init")
                    .help("Initializes a default config file")
                    .required(false)
                    .takes_value(false),
            )
            .arg(
                Arg::with_name("command")
                    .help("Runs a command that specified in config file")
                    .multiple(false)
                    .required(false)
                    .takes_value(true),
            )
            .arg(
                Arg::with_name("args")
                    .help("Sends parameters for a command")
                    .multiple(true)
                    .required(false)
                    .takes_value(true),
            )
            .get_matches();
        Self {
            arg_matches: arg_matches,
        }
    }

    pub fn run(self) -> exitcode::ExitCode {
        let code;
        let mut term = Terminal::new();

        if self.arg_matches.is_present("init") {
            match self.init_file() {
                Err(e) => {
                    code = exitcode::OSFILE;
                    term.error(&format!("{}", e));
                }
                _ => {
                    code = exitcode::OK;
                    term.write_line("Done.")
                }
            }
            return code;
        }

        match self.get_config() {
            Ok(config) => {
                code = match self.arg_matches.value_of("command") {
                    Some(command) => self.run_commands(term, config, command),
                    None => self.print_usage(term, config),
                };
            }
            Err(msg) => {
                code = exitcode::OSFILE;
                term.error(&format!("{}", msg));
            }
        }

        code
    }

    fn print_usage(self, mut term: Terminal, config: Config) -> exitcode::ExitCode {
        term.title("Environments");
        for environment in config.environments {
            term.write_line(&format!(" - {}", environment));
        }
        term.new_line();

        term.title("Commands");
        for (name, command) in config.commands {
            term.write_line(&format!(
                "  > {name:20} {help_text}",
                name = name,
                help_text = command.help_text
            ));
        }

        exitcode::OK
    }

    fn run_commands(
        &self,
        mut term: Terminal,
        config: Config,
        command: &str,
    ) -> exitcode::ExitCode {
        let environment = &env::var(ENV_NAME).unwrap_or("Unspecified".to_string());
        let args: Vec<&str> = match self.arg_matches.values_of("args") {
            Some(args) => args.collect(),
            None => vec![],
        };

        term.write(&format!("{:>11}: ", "Environment"), Some(term::color::CYAN));
        term.write(environment, None);
        term.new_line();

        match config.run_command(&command, args) {
            Ok(_) => {
                term.write(&format!("{:>11}: ", "Command"), Some(term::color::CYAN));
                term.write(&format!("{}, ", command), None);
                term.new_line();
            }
            Err(msg) => {
                term.error(&msg);
                return exitcode::DATAERR;
            }
        }
        exitcode::OK
    }

    fn init_file(&self) -> std::io::Result<()> {
        if Path::new(FILE_NAME).exists() {
            let error_msg = format!("The file {} already exists.", FILE_NAME);
            return Err(std::io::Error::new(std::io::ErrorKind::Other, error_msg));
        }
        let mut config = Config::new();
        config.get_default();

        let content = serde_yaml::to_string(&config).unwrap();
        let mut config_file = File::create(FILE_NAME)?;
        config_file.write_all(content.as_bytes())?;

        Ok(())
    }

    fn get_config(&self) -> Result<Config, ConfigError> {
        let mut current_dir = env::current_dir().unwrap();
        loop {
            if current_dir.join(FILE_NAME).exists() {
                break;
            }
            match current_dir.parent() {
                Some(parent_dir) => current_dir = parent_dir.to_path_buf(),
                None => {
                    let msg = format!("Config file {} could not be found.", FILE_NAME);
                    return Err(ConfigError::Other(msg));
                }
            }
        }

        let mut config_str = String::new();
        let mut config_file = File::open(current_dir.join(FILE_NAME)).unwrap();
        config_file.read_to_string(&mut config_str).unwrap();

        let config: Config = serde_yaml::from_str(&config_str)?;
        Ok(config)
    }
}
