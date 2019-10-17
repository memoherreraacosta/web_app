import React from 'react';
export default class Client extends React.Component {
  state = {
    books: []
  }

  componentDidMount() {
    fetch('/todos_libros/1')
      .then(res => res.json())
      .then(books => {this.setState({ books })
      console.log(this.state.books);
  });
  }

  render() {
    return (
      <ul>
        { this.state.books.map(books => 
            <li key= {books.id}>{books.titulo + "  ->  "}{books.autor}</li>
        )}
      </ul>
    )
  }
}