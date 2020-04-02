## EURO TRUCK 2 SYNC TOOL

This is repository of Euro Truck 2 sync tool, which allows to sync job list with other players.

## Hosted version
**Hosted version of this tool** [https://ft-t.github.io/](https://ft-t.github.io/)

## Changing Save format
You need to open `config.cfg` in the game folder in Documents, and change `g_developer` to 1 (one), `g_save_format`to 2 (two), and `g_console` to 1 (one). Here's a step-by-step on how to do those changes:

1. Make sure the game is closed. This won't work if you do it with the game running.
2. In the “Game Settings” panel, there is a line saying “Settings Folder: <folder with your ETS2/ATS settings> - Options”. Click “Options” → “Open game config file”. A Notepad window will appear, with the game settings file open.
3. In that Notepad window, go to “Edit” → “Find…” (or hit <kbd>Ctrl</kbd>+<kbd>F</kbd>). In the search window, type `g_developer` (notice the underscore) and hit <kbd>Enter</kbd>. It should highlight a line in the file that reads `uset g_developer "0"`. If the number isn't already `1`, change it to `1`, so that it reads `uset g_developer "1"`.
4. Go to “Edit” → “Find…” (or hit <kbd>Ctrl</kbd>+<kbd>F</kbd>) again. In the search window, type `g_save_format` and hit <kbd>Enter</kbd>. It should highlight a line in the file that reads `uset g_save_format "2"` or another number. If the number is not zero, change the number to `0`, so that it reads `uset g_save_format "0"`.
5. Go to “Edit” → “Find…” (or hit <kbd>Ctrl</kbd>+<kbd>F</kbd>) again. In the search window, type `g_console` and hit <kbd>Enter</kbd>. It should highlight a line in the file that reads `uset g_console "0"`. If the number isn't already `1`, change it to `1`, so that it reads `uset g_console "1"`.
6. Save the file and close Notepad.

## [Dev] Compiling The Source Code

 - Golang 1.14+
 - Dep

## [Dev] Hosting own server

[https://github.com/ft-t/ets2-sync/blob/master/.deploy/app-playbook.yml](https://github.com/ft-t/ets2-sync/blob/master/.deploy/app-playbook.yml)
