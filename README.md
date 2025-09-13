# Validate GPU‑param fitting with Ebiten

## Goal
Validate in ~14 hours whether a small fragment‑shader pipeline (Ebiten/Kage) plus CPU loss (L1 and optionally SSIM) can fit parameters that visibly reduce error on a few image pairs, with a live preview window for debugging.

## Day 1 (6–7h): Minimal pipeline + L1 loss + live preview
- [ ] Repo scaffold (cmd tool + internal packages: pipeline, loss, opt, io)
- [ ] Load N small image pairs (inputs/, targets/), apply EXIF orientation, downsample to 256–512 px
- [ ] Build offscreen render pipeline (ebiten.Images ping‑pong)
 - [ ] Passes: exposure → 1D tone curve (texture) → unsharp/saturation (keep it to 2–3 passes)
 - [ ] Linearize at start, convert to sRGB only on final write
- [ ] ReadPixels once per evaluation to a preallocated buffer
- [ ] CPU L1 loss over RGBA or luma (pick one and stick to it)
- [ ] Simple optimizer loop (random search + Nelder–Mead via Gonum, small eval budget)
- [ ] Live preview window
 - [ ] Draw latest output and/or best‑so‑far; side‑by‑side with target and a diff heatmap
 - [ ] Throttle UI to 10–30 FPS; update text (iter, loss) every ~200–500 ms
- [ ] CLI flags: input dir, target dir, max size, evals, seed
- [ ] Log best θ and save model.json; add “apply” mode to render full‑res outputs for 1–2 images

## Day 2 (7–8h): SSIM, speedups, validation
- [ ] Add single‑scale SSIM (windowed 11×11, Gaussian) on CPU; combine loss = α·L1 + β·(1−SSIM)
- [ ] Parallelize CPU loss across image pairs; keep GPU work on Ebiten’s thread
- [ ] Reduce overhead
 - [ ] Reuse ebiten.Images between evals
 - [ ] Keep one ReadPixels buffer; avoid allocations in hot loop
 - [ ] Optionally shrink eval size (e.g., 384 px) if slow
- [ ] Optimizer tightening
 - [ ] Budgeted restarts; early stop if no improvement in K iters
 - [ ] Save checkpoints of best θ and preview PNG
- [ ] Validation
 - [ ] Plot loss vs. iteration (CSV + quick script) and eyeball improvement
 - [ ] Run “apply” on 3–5 full‑res images; visually inspect artifacts
- [ ] Debrief: note bottlenecks (readback, SSIM time) and gaps (need GPU reduction?) to decide next steps

## Things to watch out for
- sRGB vs linear: do math in linear; only encode to sRGB at the very end
- Readbacks: the main cost; keep images small and one readback per eval
- Threading: all Ebiten GPU calls on its main thread; do CPU loss/optimizer in goroutines
- Determinism: fix RNG seed; keep image loading and pipelines stable across runs
- Precision/banding: 8‑bit outputs; consider subtle dithering on the final pass if banding appears
- Parameter scaling: normalize θ ranges so the optimizer behaves (e.g., map UI ranges to [−1,1])

## Handy links
- Ebiten site and docs: https://ebiten.org/
- Ebiten API (pkg.go.dev): https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2
 - Image.ReadPixels: https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.ReadPixels
 - SetMaxTPS: https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#SetMaxTPS
 - Kage shader docs: https://ebiten.org/documents/shader.html
 - Examples: https://github.com/hajimehoshi/ebiten/tree/main/examples
- Go image resampling: https://pkg.go.dev/golang.org/x/image/draw
- Gonum optimize (Nelder–Mead): https://pkg.go.dev/gonum.org/v1/gonum/optimize
- SSIM reference (overview): https://en.wikipedia.org/wiki/Structural_similarity
