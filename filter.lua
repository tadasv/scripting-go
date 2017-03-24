function filter(line)
    found = string.find(line, "if")
    if found == nil then
        return false
    end

    return true
end
