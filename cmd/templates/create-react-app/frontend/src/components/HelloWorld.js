import React from 'react';
import Modal from 'react-modal';

class HelloWorld extends React.Component {
  constructor(props, context) {
    super();
    this.state = {
      showModal: false
    };

    this.handleOpenModal = this.handleOpenModal.bind(this);
    this.handleCloseModal = this.handleCloseModal.bind(this);
  }

  handleOpenModal () {
    this.setState({ showModal: true });

    window.backend.basic().then(result =>
      this.setState({
        result
      })
    );
  }

  handleCloseModal () {
    this.setState({ showModal: false });
  }

  render() {
    const { result } = this.state;
    return (
      <div className="App">
        <button onClick={this.handleOpenModal} type="button">
          Hello
        </button>
        <Modal
          isOpen={this.state.showModal}
          contentLabel="Minimal Modal Example"
        >
        <p>{result}</p>
        <button onClick={this.handleCloseModal}>Close Modal</button>
        </Modal>
      </div>
    );
  }
}

export default HelloWorld;
