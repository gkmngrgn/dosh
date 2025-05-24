-- I've a config repo, I configure my system using this script.

local config_dir = "~/.config"

cmd.add_task{
   name="config_os",
   description="copy my configuration files and replace",
   command=function()
      -- copy config files.
      cmd.copy("./config/*", config_dir)
      cmd.copy("./home/.*", "~")

      -- zsh specific settings
      if env.IS_ZSH then
         if not cmd.exists("~/.oh-my-zsh") then
            cmd.run_url("https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh")
         end

         if cmd.exists_command("conda") then
            cmd.run("conda init zsh")
         end
      end

      -- tmux
      cmd.clone{
         url="https://github.com/tmux-plugins/tpm",
         destination="~/.tmux/plugins/tmp",
         fetch=true,
      }
   end
}

cmd.add_task{
   name="install_cli_apps",
   description="install my favourite apps",
   command=function ()
      if env.IS_WINDOWNS then
         local packages = {"Git.Git", "VSCodium.VSCodium", "Discord.Discord", "Valve.Steam"}
         cmd.winget_install(packages)
      elseif env.IS_MACOS then
         local packages = {
            "MisterTea/et/et", "bat", "clojure", "cmake", "deno", "exa", "exercism", "fd",
            "git-delta", "git-lfs", "golang", "helix", "htop", "hugo", "jq", "llvm",
            "multimarkdown", "openssl", "pass", "pre-commit", "ripgrep", "rustup-init",
            "rust-analyzer", "shellcheck", "tmux", "font-ibm-plex", "miktex-console", "miniconda"
         }
         local taps = {"helix-editor/helix"}
         cmd.brew_install{packages, cask=true, taps=taps}
      elseif env.IS_LINUX then
         local packages = {"git", "ripgrep"}
         cmd.apt_install(packages)
      end

      if not cmd.IS_WINDOWS and not cmd.exists_command("nvm") then
         cmd.run_url("https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh")
      end
   end
}

cmd.add_task{
   name="change_theme",
   description="sync your editor theme with system",
   command=function (name)
      if name == "dark" then
         cmd.info("modus vivendi.")
      elseif name == "light" then
         cmd.info("modus operandi.")
      else
         cmd.error("invalid profile. choose `light` or `dark`.")
      end
   end
}
