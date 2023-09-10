# Pipni API

Pipni API is a desktop application for testing HTTP RESTful endpoints.

![screenshot](screenshot.png)

**WARNING:** This project is still in its infancy and there is a lot to be fixed and added. You will find it hard to achieve your day-to-day goals with the current state.


## Getting Started

Stable versions of the project will be available in the Releases section of the repository.

Alternatively, if you have Go properly setup, you can fetch the project as follows:

```sh
git clone https://github.com/mokiat/PipniAPI
cd PipniAPI
```

And you can run it as follows.

```sh
go run ./cmd/PipniAPI
```

Though the preferred way is to use the Taskfile.

```sh
task run
```

## Goal and Purpose

There are alternative tools out there for making REST calls through a user interface. However, those projects are either payed or are on their way to becoming payed, are bloated, and gradually moving to the cloud.

The plan for this project is to always be free and open-source. Furthermore, it will always remain a pure desktop app (no cloud syncing and accounts required). Lastly, it aims to remain fairly simple - just REST API calls or similar - no SOAP, gRPC or what-not.

If for some reason this project ever decides to deviate from the above, it will be forked with a separate name and this repository will remain true to its original goals.


## Contributions

The project is still very early in its development stages. The code structure and design patterns may need to be changed. What's more, the UI is based on the [lacking](https://github.com/mokiat/lacking) game framework's UI package. It was never intended for such a project but for a number of reasons it was the one that was picked for this project (more information below).

As such, for now the best way to contribute is to open an Issue.

## Why in Go and why a custom UI framework?

What follows is my personal opinion on the matter so don't take it personal if it strikes too close to the heart.

First, I am a fan of the Go programming language. I find it very easy to write. It is very readable. And complex scenarios take little effort to implement correctly. In 90% of the time it is as fast as one needs it to be and in 5% one can come up with workarounds (the built-in benchmarking helps a lot). The last 5% are rare. As such, there is currently no other language that would bring me as much joy (and motivation) to do this in.

That said, I did try ElectronJS in combination with Svelte and Tailwind CSS. Getting these libraries to work together was a mess. There seem to be multiple frameworks and tools to get a NodeJS project working (transpiling, packaging, relocting, watching) these days and it was not clear till the very end whether I had configured the correct thing and whether I was up to date or already legacy in my approach. By the end of the bootstrapping, my energy had already vanished. This is why I opted for the simple `go` command - always the same, everywhere. I wasn't looking forward having to deal with ElectronJS's sandbox rpc communication paradigm and always wondering whether I am secure or not.

Having decided on Go, I did give [Fyne](https://github.com/fyne-io/fyne) a try, since this appears to be the most popular UI framework out there. While I am very happy to see that there are efforts in this direction (since Go has been behind from its start), there were a number of things in Fyne that didn't work out well for me. First, the UI approach is imperative, whereas I have grown to prefer the declarative approach of frameworks like Svelte, Vue, and ReactJS. Second, there were some inconsistencies in the API - some things were modified through setters and others were modified through public fields (having to call Refresh afterwards). Lastly, I did not like the concurrent approach of the API. I strongly believe that a UI framework should be single-threaded with mechanisms to coordinate with background jobs, otherwise the framework becomes difficult to write with potential for bugs and the client code needs to do all the synchronization itself, when in fact 99% of the code would have worked best on a single go routine.

Lastly, I had already started work on the UI framework of lacking a very long time ago. I still need it for some toy game projects, so I might as well approach it from multiple angles. I am aware that there is a lot missing but this is a risk (and a challenge) I am willing to take (embark on).
