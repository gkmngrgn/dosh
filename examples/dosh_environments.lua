-- limit command use by the environment and operating system

local cmd = require("dosh_commands")

cmd.add_task{
   name="train_data",
   description="train data in server",
   environments={ "prod", "stag" },
   required_platforms={ "linux" },
   command=function ()
      cmd.run_url("http://localhost:8000/train_data.sh")
   end
}
