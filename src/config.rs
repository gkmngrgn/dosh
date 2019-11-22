use serde::{Deserialize, Serialize};
use serde_yaml::{Mapping, Number, Value};
use std::collections::HashMap;
use std::fmt;

pub static FILE_NAME: &str = "dosh.yaml";

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Config {
    settings: Mapping,
    pub environments: Vec<String>,
    aliases: HashMap<String, Alias>,
    pub commands: HashMap<String, Command>,
}

impl Config {
    pub fn new() -> Self {
        Self {
            settings: Mapping::new(),
            environments: vec!["PROD".to_string(), "TEST".to_string()],
            aliases: HashMap::new(),
            commands: HashMap::new(),
        }
    }

    pub fn get_default(&mut self) {
        // A default config for managing hugo website.
        self.settings.insert(
            Value::String("version".to_string()),
            Value::String(env!("CARGO_PKG_VERSION").to_string()),
        );
        self.settings.insert(
            Value::String("verbosity".to_string()),
            Value::Number(Number::from(0)),
        );

        self.aliases.insert(
            "hugo".to_string(),
            Alias {
                default: "hugo".to_string(),
                linux: None,
                mac: None,
                win: Some("hugo.exe".to_string()),
            },
        );

        self.commands.insert(
            "start".to_string(),
            Command {
                environments: vec![],
                help_text: "run the development server".to_string(),
                run: CommandRun::OneLine(String::from("{hugo} server -D")),
            },
        );
        self.commands.insert(
            "build".to_string(),
            Command {
                environments: vec![],
                help_text: "generate static files to the public folder".to_string(),
                run: CommandRun::OneLine(String::from("{hugo}")),
            },
        );
        self.commands.insert(
            "deploy".to_string(),
            Command {
                environments: vec!["PROD".to_string()],
                help_text: "deploy to the production server".to_string(),
                run: CommandRun::Group(vec![
                    String::from("CMD build"),
                    String::from("CD public"),
                    String::from("RUN git add ."),
                    String::from("RUN git commit -m \"publish changes.\""),
                    String::from("RUN git push origin master"),
                    String::from("CD .."),
                    String::from("PRINT \"DONE.\""),
                ]),
            },
        );
    }

    pub fn run_command(&self, command: &str, _: Vec<&str>) -> Result<(), String> {
        // TODO: second parameter is for the command arguments. Name it as args later.
        match self.commands.get(command) {
            Some(cmd) => {
                // TODO: check env here.
                let cmd_runs: Vec<String> = match &cmd.run {
                    CommandRun::OneLine(run) => vec![run.to_string()],
                    CommandRun::Group(runs) => runs.clone(),
                };
                for cmd_run in cmd_runs {
                    println!("{:?}", cmd_run);
                }
            }
            None => return Err(format!("Unknown command: {}.", command)),
        }
        Ok(())
    }
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Alias {
    default: String,
    linux: Option<String>,
    mac: Option<String>,
    win: Option<String>,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Command {
    pub help_text: String,
    run: CommandRun,
    environments: Vec<String>,
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
#[serde(untagged)]
pub enum CommandRun {
    OneLine(String),
    Group(Vec<String>),
}

#[derive(Debug)]
pub enum ConfigError {
    Yaml(serde_yaml::Error),
    Other(String),
}

impl From<serde_yaml::Error> for ConfigError {
    fn from(error: serde_yaml::Error) -> Self {
        ConfigError::Yaml(error)
    }
}

impl fmt::Display for ConfigError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            ConfigError::Yaml(ref e) => write!(f, "{}", e),
            ConfigError::Other(ref e) => write!(f, "{}", e),
        }
    }
}
