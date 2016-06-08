import * as CodeMirror from 'codemirror'
import * as $ from 'jquery'
import * as moment from 'moment';
import * as React from './dom';

import GitHub from './github';

import 'codemirror/mode/ruby/ruby';
import 'codemirror/addon/runmode/runmode';

export default class Editor {
    editor: CodeMirror.Editor;
    gh: GitHub;
    lobbyId: string;

    constructor(lobbyId: string, token: string) {
        this.gh = new GitHub(token)
        this.lobbyId = lobbyId

        let pane = document.getElementById('editor-pane')
        this.editor = CodeMirror(pane, {
            lineNumbers: true,
            mode: 'ruby',
        })

        this.editor.refresh()
        document.getElementById('editor-modal').addEventListener("resize", () => {
            this.editor.refresh()
            this.editor.focus()
        })

        $('#editor-run-btn').click(() => {
            this.send(() => this.close())
        })

        $('#editor-load-btn').click(() => {
            this.updateGists();
            $('#gist-pane').addClass('showPane')
        })

        $('#gist-cancel-btn').click(() => {
            $('#gist-pane').removeClass('showPane')
        })

        let aiUrl = "/lobbies/" + lobbyId + "/ai";
        $.ajax({
            url: aiUrl,
        }).done((data) => {
            this.editor.setValue(data)
        })

        this.updateGists();
    }

    private async updateGists() {
        let gists = await this.gh.getGists();
        let items = [];

        for (let gist of gists) {
          let file = gist.files[Object.keys(gist.files)[0]];
        console.log(gist);

          if (file.language != "Ruby") {
            continue;
          }

          let contents = await $.ajax({ url: file.raw_url });
          let code = <pre class="cm-s-default gist-snippet" />;
          CodeMirror.runMode(contents, "ruby", code);

          let name = gist.description || file.filename;
          let updated = moment(gist.updated_at).fromNow();
          let el =
            <div class="gist-entry">
              <div class="gist-info">
                <span class="gist-name">{name}</span>
                <span class="gist-date">Updated {updated}</span>
              </div>
              {code}
            </div>;

          items.push(el);
          items.push(<hr/>)
        }

        let list = document.getElementById('gist-list')
        while (list.firstChild) list.removeChild(list.firstChild);
        for (let item of items) {
          list.appendChild(item)
        }
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

    public send(cb) {
        let aiUrl = "/lobbies/" + this.lobbyId + "/ai";
        let data = this.editor.getValue();
        $.post(aiUrl, data, function() { cb() })
    }
}
