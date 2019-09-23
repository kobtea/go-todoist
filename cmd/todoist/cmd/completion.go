package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	bashCompletionFunc = `
__todoist_select_one() {
	fzf
}

__todoist_select_multi() {
	fzf -m
}

__todoist_filter_ids() {
	COMPREPLY=( $(todoist filter list | __todoist_select_multi | awk '{print $1}' | tr '\n' ' ') )
}

__todoist_item_id() {
	COMPREPLY=( $(todoist item list | __todoist_select_one | awk '{print $1}') )
}

__todoist_item_ids() {
	COMPREPLY=( $(todoist item list | __todoist_select_multi | awk '{print $1}' | tr '\n' ' ') )
}

__todoist_label_id() {
	COMPREPLY=( $(todoist label list | __todoist_select_one | awk '{print $1}') )
}

__todoist_labels_ids() {
	COMPREPLY=( $(todoist label list | __todoist_select_multi | awk '{print $1}' | tr '\n' ' ') )
}

__todoist_project_id() {
	COMPREPLY=( $(todoist project list | __todoist_select_one | awk '{print $1}') )
}

__todoist_project_ids() {
	COMPREPLY=( $(todoist project list | __todoist_select_multi | awk '{print $1}' | tr '\n' ' ') )
}

__todoist_custom_func() {
	case ${last_command} in
		todoist_filter_update | todoist_filter_delete)
			__todoist_filter_ids
			return
			;;
		todoist_item_update | todoist_item_delete | todoist_item_move | todoist_item_complete | todoist_item_uncomplete)
			__todoist_item_ids
			return
			;;
		todoist_label_update | todoist_label_delete)
			__todoist_label_ids
			return
			;;
		todoist_project_update | todoist_project_delete | todoist_project_archive | todoist_project_unarchive)
			__todoist_project_id
			return
			;;
		*)
			;;
	esac
}
`
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate completion script",
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "generate bash completion script",
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

// Refs
// - https://github.com/kubernetes/kubernetes/blob/master/pkg/kubectl/cmd/completion.go
// - https://github.com/spf13/cobra/issues/107#issuecomment-270140429
var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "generate zsh completion script",
	Run: func(cmd *cobra.Command, args []string) {
		zshHead := "#compdef todoist"
		zshInitialization := `
__todoist_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}

__todoist_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__todoist_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__todoist_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__todoist_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__todoist_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__todoist_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__todoist_filedir() {
	local RET OLD_IFS w qw
	__todoist_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __todoist_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__todoist_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__todoist_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
		printf %q "$1"
    fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__todoist_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__todoist_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__todoist_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__todoist_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__todoist_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__todoist_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__todoist_type/g" \
	<<'BASH_COMPLETION_EOF'
`
		zshTail := `
BASH_COMPLETION_EOF
}
__todoist_bash_source <(__todoist_convert_bash_to_zsh)
_complete todoist 2>/dev/null
`
		fmt.Println(zshHead)
		fmt.Print(zshInitialization)
		RootCmd.GenBashCompletion(os.Stdout)
		fmt.Print(zshTail)
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
}
