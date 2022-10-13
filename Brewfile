# This Brewfile uses dependencies defined in .goreleaser.yaml, so they are kept in sync
# https://goreleaser.com/customization/homebrew/#:~:text=package%20depends%20on.-,dependencies,-%3A%0A%20%20%20%20%20%20%2D%20name

install_optional = ENV.fetch('INSTALL_OPTIONAL', 'true') == 'true'

YAML.load_file('.goreleaser.yaml')['brews'].each do |tap|
  tap['dependencies'].each do |dep|
    next unless install_optional

    if dep['version']
      brew "#{dep['name']}@#{dep['version']}"
    else
      brew dep['name']
    end
  end
end
