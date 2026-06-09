# .bashrc — managed by GNU Stow
# Sources Omarchy defaults, then modular .bashrc.d/ overrides

shopt -s checkwinsize histappend autocd cdspell dirspell globstar

HISTFILESIZE=200000
HISTSIZE=100000
HISTCONTROL=ignoreboth:erasedups
HISTTIMEFORMAT='%F %T '

# Omarchy defaults (aliases, env, functions, etc.)
source ~/.local/share/omarchy/default/bash/rc

# Local overrides
if [ -d "$HOME/.bashrc.d" ]; then
    for i in "$HOME/.bashrc.d"/*.sh; do
        [ -r "$i" ] && . "$i"
    done
fi
unset i
