export function lerp(from, to, amount) {
    return from * (1 - amount) + to * amount
}
