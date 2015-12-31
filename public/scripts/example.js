var Player = React.createClass({
  rawMarkup: function() {
    var rawMarkup = marked(this.props.children.toString(), {sanitize: true});
    return { __html: rawMarkup };
  },

  render: function() {
    return (
      <div className="player">
        <h2 className="playerName">
          {this.props.name}
        </h2>
        <span dangerouslySetInnerHTML={this.rawMarkup()} />
      </div>
    );
  }
});

var PlayerBox = React.createClass({
  loadPlayersFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    this.loadPlayersFromServer();
    setInterval(this.loadPlayersFromServer, this.props.pollInterval);
  },
  render: function() {
    return (
      <div className="playerBox">
        <h1>Players</h1>
        <PlayerList data={this.state.data} />
      </div>
    );
  }
});

var PlayerList = React.createClass({
  render: function() {
    var playerNodes = this.props.data.map(function(player) {
      return (
        <Player name={player.name} key={player.id}>
          {player.quality}
        </Player>
      );
    });
    return (
      <div className="playerList">
        {playerNodes}
      </div>
    );
  }
});


ReactDOM.render(
  <PlayerBox url="/api/players" pollInterval={2000} />,
  document.getElementById('content')
);
