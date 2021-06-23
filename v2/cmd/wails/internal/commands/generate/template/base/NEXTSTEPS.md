# Next Steps

Congratulations on generating your template!

## Completing your template

The next steps to complete the template are:

  1. Complete the fields in the `template.json` file. 
  2. Update `README.md`.
  3. Edit `wails.tmpl.json` and ensure all fields are correct, especially:
     - `html` - path to your `index.html`
     - `frontend:install` - The command to install your frontend dependencies
     - `frontend:build` - The command to build your frontend
  4. Delete this file.

## Testing your template

You can test your template by running this command:

`wails init -name test -template /path/to/template`

### Checklist 
Once generated, do the following tests:
  - Change into the new project directory and run `wails build`. A working binary should be generated in the `build/bin` project directory.
  - Run `wails dev`. This will compile your backend and run it. You should be able to go into the frontend directory and run `npm run dev` (or whatever your dev command is) and this should run correctly. You should be able to then open a browser to your local dev server and the application should work.