import * as CodeMirror from 'codemirror'
import * as $ from 'jquery'
import * as moment from 'moment';
import * as React from './dom';
import * as sweetalert from 'sweetalert'

import GitHub from './github';

import 'codemirror/mode/ruby/ruby';
import 'codemirror/addon/runmode/runmode';

export default class Editor {
    editor: CodeMirror.Editor;
    gh: GitHub;
    currentGist: any;
    lobbyId: string;

    constructor(lobbyId: string, token: string) {
        this.gh = new GitHub(token)
        this.lobbyId = lobbyId
        this.currentGist = null;

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
            this.updateGistList();
            $('#gist-pane').addClass('showPane')
        })

        $('#editor-create-gist-btn').click(() => {
            $('#save-pane').addClass('showPane')
        })

        $('#gist-cancel-btn').click(() => {
            $('#gist-pane').removeClass('showPane')
        })

        $('#editor-save-gist-btn').click(() => {
          let cb = (gist) => {
            this.currentGist = gist;
            swal("Saved !", "Gist saved !", "success");
            this.updateGistList()
          }

          if (this.currentGist !== null) {
            this.updateGist(cb)
          } else {
            this.createGist(cb)
          }
        })

        let aiUrl = "/lobbies/" + lobbyId + "/ai";
        $.ajax({
            url: aiUrl,
        }).done((data) => {
            this.editor.setValue(data)
        })

        this.updateGistList();
    }

    private gistElement(gist) {
      let file = gist.files[Object.keys(gist.files)[0]];

      if (file.language != "Ruby") {
        return null;
      }

      return $.ajax({ url: file.raw_url }).then((contents) => {
        let code = <pre class="cm-s-default gist-snippet" />;
        CodeMirror.runMode(contents, "ruby", code);

        let name = gist.description || file.filename;
        let updated = moment(gist.updated_at).fromNow();
        let el =
          <div class="gist-entry"
               onClick={() => {
                 this.currentGist = gist;
                 this.editor.setValue(contents)
                 $('#gist-pane').removeClass('showPane')
               }}>

            <div class="gist-info">
              <span class="gist-name">{name}</span>
              <span class="gist-date">Updated {updated}</span>
            </div>
            {code}

          </div>;
          return el;
      })
    }

    private updateGistList() {
        this.gh.getGists().then((gists) => {
          let elements = gists.map((gist) => this.gistElement(gist));
          return $.when(...elements);
        }).then((...items) => {
            let list = document.getElementById('gist-list')
            while (list.firstChild) list.removeChild(list.firstChild);
            items.forEach((item) => {
              if (item !== null) {
                list.appendChild(item as any)
              }
            })
        })
    }

    private updateGist(cb) {
      let contents = this.editor.getValue()
      let filename = Object.keys(this.currentGist.files)[0]
      let files = {}
      files[filename] = {
        content: contents,
        language: "Ruby",
      }

      this.gh.updateGist(this.currentGist.id, {
          files: files,
      }).then(cb)
    }

    private createGist(cb) {
      let contents = this.editor.getValue()

      sweetalert({
        title: "Save as Gist",
        text: "Description :",
        type: "input",
        showCancelButton: true,
        closeOnConfirm: false,
        animation: "slide-from-top",
      }, (desc) => {
        if (desc === false) return false;
        desc = (desc as string).trim()
        if (desc === "") {
          return false
        }

        this.gh.createGist({
          description: desc,
          files: {
            "ai.rb": {
              content: contents,
              language: "Ruby",
            }
          },
        }).then(cb)
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

    public send(cb) {
        let aiUrl = "/lobbies/" + this.lobbyId + "/ai";
        let data = this.editor.getValue();
        $.post(aiUrl, data, function() { cb() })
    }
}
