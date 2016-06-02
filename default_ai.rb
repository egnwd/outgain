angle = rand * 2 * Math::PI
@dx = Math::cos(angle)
@dy = Math::sin(angle)

def run(player, world)
    angle = rand * Math::PI / 2 - Math::PI / 4

    c = Math::cos(angle)
    s = Math::sin(angle)

    dx = c * @dx - s * @dy
    dy = s * @dx + c * @dy

    @dx = dx
    @dy = dy

    target = world.entities
        .select {|e| e.id != player.id and e.radius < player.radius }
        .min_by {|e| player.distance e}

    if target != nil then
        return [target.x - player.x, target.y - player.y]
    else
        return [@dx, @dy]
    end
end
