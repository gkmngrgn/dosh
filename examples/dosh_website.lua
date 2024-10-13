-- this is how I deploy my hugo websites.

local function check_commands (commands)
   for i = 1, #commands do
      if not cmd.exists_command(commands[i]) then
         error("Command does not exist: " .. commands[i])
      end
   end
end

cmd.add_task{
   name="deploy",
   description="deploy changes",
   command=function ()
      check_commands({"hugo", "rsync"})
      cmd.run("hugo")
      cmd.run("rsync -avz --delete public/ user@mydomain.com:/var/www/html/mydomain.com/")
   end
}

cmd.add_task{
   name="dev",
   description="run development server",
   command=function ()
      check_commands({"hugo"})
      cmd.run("hugo server -D --bind 0.0.0.0 -b myngrokdomain")
   end
}
