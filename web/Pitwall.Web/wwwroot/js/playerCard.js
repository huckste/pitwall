// Tilts the player card toward the pointer and moves its foil/glare with it.
// Kept small and snappy so it reads as a rigid card, not something floating.
// Runs in plain JS so it stays smooth (no .NET round-trip per mouse move).
const MAX_TILT_X = 6; // degrees, top/bottom
const MAX_TILT_Y = 8; // degrees, left/right

export function attachTilt(card) {
  const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const reset = () => {
    card.classList.remove("tracking");
    card.style.setProperty("--mx", "0.5");
    card.style.setProperty("--my", "0.5");
    card.style.setProperty("--rx", "0deg");
    card.style.setProperty("--ry", "0deg");
  };

  const move = (e) => {
    const r = card.getBoundingClientRect();
    const x = Math.min(Math.max((e.clientX - r.left) / r.width, 0), 1);
    const y = Math.min(Math.max((e.clientY - r.top) / r.height, 0), 1);
    card.classList.add("tracking");
    card.style.setProperty("--mx", x.toFixed(3));
    card.style.setProperty("--my", y.toFixed(3));
    if (!reduceMotion) {
      card.style.setProperty("--rx", `${((0.5 - y) * 2 * MAX_TILT_X).toFixed(2)}deg`);
      card.style.setProperty("--ry", `${((x - 0.5) * 2 * MAX_TILT_Y).toFixed(2)}deg`);
    }
  };

  reset();
  card.addEventListener("pointermove", move);
  card.addEventListener("pointerleave", reset);

  return {
    dispose() {
      card.removeEventListener("pointermove", move);
      card.removeEventListener("pointerleave", reset);
    },
  };
}
