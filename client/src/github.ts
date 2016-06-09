import * as $ from 'jquery'

export default class Github {
  token: string;
  base: string;

  constructor(token: string, base?: string) {
    this.token = token
    this.base = base || 'https://api.github.com';
  }

  private request(method: string, path: string, data?: any) {
    if (data !== null) {
      data = JSON.stringify(data)
    }
    return $.ajax({
      method: method,
      url: this.base + path,
      headers: {
        Accept: "application/vnd.github.v3+json",
        Authorization: "token " + this.token,
      },
      dataType: 'json',
      data: data,
    })
  }

  public getGists() {
    return this.request("GET", "/gists");
  }

  public createGist(gist) {
    return this.request("POST", "/gists", gist);
  }

  public updateGist(id, gist) {
    return this.request("PATCH", "/gists/" + id, gist);
  }
}
