'use strict';

class ArticleRow extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <tr>
        <td>{this.props.article.first_seen}</td>
        <td>JVECC</td>
        <td><a href={this.props.article.url}>{this.props.article.title}</a></td>
      </tr>
    );
  }
}

class ArticleList extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      articles: data.articles
    };
  }

  render() {
    var rows = this.state.articles.map((a) =>
      <ArticleRow article={a} />
    );
    return (
      <table>
        <thead>
          <th>
            <td>Date</td>
            <td>Journal</td>
            <td>Title</td>
          </th>
        </thead>
        <tbody>
          {rows}
        </tbody>
      </table>
    );
  }
}

const domContainer = document.querySelector('#article-list');
ReactDOM.render(React.createElement(ArticleList), domContainer);
