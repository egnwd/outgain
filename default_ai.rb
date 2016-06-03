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

    prey = world.entities
        .select {|e| e.resource? or (e.creature? and e.id != player.id and e.radius < player.radius) }
        .min_by {|e| player.distance e}

    predator = world.entities
        .select {|e| e.spike? or (e.creature? and e.id != player.id and e.radius > player.radius) }
        .min_by {|e| player.distance e}

    if prey != nil then
        return [prey.x - player.x, prey.y - player.y]
    elsif predator != nil then
        return [player.x - predator.x, player.y - predator.y]
    else
        return [@dx, @dy]
    end
end
