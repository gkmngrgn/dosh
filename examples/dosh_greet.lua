-- say hello to anyone. it takes an argument.
cmd.add_task{
   name="say_hello",
   description="say hello to anyone",
   command=function(there)
      there = there or "world"
      cmd.info("hello " .. there .. "!")
   end
}

-- say your current OS.
cmd.add_task{
   name="whatismyos",
   description="print my current operating system",
   command=function()
      local message

      if env.IS_MACOS then
         message = "Your OS is macOS"
      elseif env.IS_WINDOWS then
         message = "Your OS is Windows"
      elseif env.IS_LINUX then
         message = "Your OS is Linux"
      else
         message = "What the hell are you using?"
      end

      cmd.run("echo '" .. message .. "'")
   end
}
