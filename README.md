# ShadyEdits

ShadyEdits is a GPU-accelerated CLI tool that tries to optimize shader parameters to minimize error between input and target images. It uses [Ebiten](https://ebiten.org/) for live preview and ebitens own shader language [Kage](https://ebiten.org/documents/shader.html).

## Demo

![Demo Video](demo.gif)

## Why

This is just a fun little side project to learn a little more about shaders and parameter optimizations.
I am sure there are faster ways to find how an image is edited.

## Supported Shaders

- **Alpha**: Adjusts image opacity and blending.
- **Exposure**: Modifies image exposure in stops.
- **Contrast**: Changes image contrast in linear RGB.
- **Saturation**: Boosts or reduces color saturation.
- **Temperature**: Warms or cools image colors.

## Supported Tuners

- **RandomSearch**: Simple random search
- **RandomGeneticEvolve**: Genetic algorithm-based parameter search.

## Installation & Usage

### Prerequisites

- Go 1.24+
- Ebiten v2
- 2 images of the same size

### Install

```sh
git clone https://github.com/seppedelanghe/shady-edits.git
cd ShadyEdits
go build ./cmd/cli
```

### Run

```sh
./cli assets/input.png assets/target.png
```

This will launch the optimizer and preview window. You can adjust input/target paths as needed.


## References

- [Ebiten](https://ebiten.org/)
- [Kage Shader Docs](https://ebiten.org/documents/shader.html)
