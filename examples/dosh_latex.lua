-- all my cover letters, resumes, presentations are written in LaTeX.
-- I use this script for generating PDF files.

cmd.add_task{
   name="build",
   description="run formatter and generate pdf files",
   required_commands={ "latexindent", "latexmk" },
   command=function ()
      cmd.run("latexindent -w -s tex/*.tex")
      cmd.run("latexmk -xelatex tex/*.tex")
      cmd.run("mkdir -p pdf/ && mv *.pdf pdf/")
   end
}

cmd.add_task{
   name="clean",
   description="remove build files, except pdfs",
   command=function ()
      cmd.run("rm *.{aux,fdb_latexmk,fls,log,out,xdv}")
      cmd.run("rm tex/*.bak*")
   end
}
