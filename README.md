# space-go
This project is to test Ebiten and try a way to deploy to all follow plataforms:
* Windows
* Linux
* Web (WASM)
* Android(_Maybe_)

## GamePlay
At the end of this project the player should be able to move a spaceship and shot enemies that will be spawned for all 4 directions randomally and score the maximo as they can.

### Game Things
* Start Menu
* Arrow Movement
* Spacebar Shoting
* Gamepad Equivalent Controller(_Maybe_)
* Screen Equivalent Controller(_Maybe_)
* Game Over Menu showing the current highest score

## Deploy
All builds are supposed to be deployed to [Itch.io Space Go](https://maiconspas.itch.io/space-go)

### Artifacts
All builds will be avaliable to download as artifacts inside the GitAction CD

## Automated test
Not defined yet

## Links used for this build

### To build for Android
[Ebiten Guide](https://ebitengine.org/en/documents/mobile.html)

[Install Android SDK](https://gist.github.com/steveclarke/d988d89e8cdf51a8a5766d69ecb07e7b)
* Extra command to install ndk: `sdkmanager ndk-bundle`
