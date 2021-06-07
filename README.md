# Main goal

This application's main goal was to test cross platform support of golang using simple code. As the lover of the old dos games abandonware i've tried to write a cool frontend for dosbox dos emulator.

# What is LegacyBest ?

Well in short it's a wrapper for dosbox with few nice features. The main goal as you know was to write cross platform application using Golang. The second goal was to connect it to the external API, it can fetch game lists with images/description and categories without any local data. Then you can install/remove games in the way you like. You can also search for game titles. The API of this application will also be shared on github soon.

# How that crossplatform compatibility is made ?

I'm using webview component, which displays bundled in application html, turns application into small webserver that serves html and listens for the events, it also launches a client to browse that content. It works on every big platform as native application, just uses web technologies, it's lightweight and does not need a lot of resources. 

# What host platforms does it support ?

* Apple Macintosh (Intel x64)
* Microsoft Windows (Intel x86/x64)
* Linus Torvalds Linux (Intel x86/x64)

# Screenshots

[](/pics/Screenshot2020-09-03at10.01.43.png)