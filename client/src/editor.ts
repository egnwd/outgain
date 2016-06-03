import * as CodeMirror from 'codemirror'
import * as $ from 'jquery'

declare function require(string): any;
require('codemirror/mode/ruby/ruby');

export class Editor {
    editor: CodeMirror.Editor;
    lobbyId: string;

    constructor(lobbyId: string) {
        this.lobbyId = lobbyId;
        let input = document.getElementById('editor-input');
        this.editor = CodeMirror(input, {
            lineNumbers: true,
            mode: 'ruby',
        });

        this.editor.refresh()
        document.getElementById('editor-modal').addEventListener("resize", () => {
            this.editor.refresh()
            this.editor.focus()
        })

        $('#editor-save-btn').click(() => {
            this.close()
        })

        $('#editor-cancel-btn').click(() => {
            this.close()
        })

        let aiUrl = "/lobbies/" + lobbyId + "/ai";
        $.ajax({
            url: aiUrl,
        })
        .done((data) => {
            this.editor.setValue(data)
        })
    }

    public open() {
        $('#editor').removeClass('hideEditor')
        $('#editor').addClass('showEditor')
        $('#editor').show()

        setTimeout(() => {
            this.editor.refresh()
            this.editor.focus()
        }, 300);
    }

    public close() {
        $('#editor').removeClass('showEditor')
        $('#editor').addClass('hideEditor')

        setTimeout(() => {
            $('#editor').hide()
        }, 300);
    }
}
