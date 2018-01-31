import React from 'react';
import ReactDOM from 'react-dom';
import GraphiQL from 'graphiql';
import fetch from 'isomorphic-fetch';
import 'graphiql/graphiql.css';

function graphQLFetcher(graphQLParams) {
  return fetch('/query', {
    method: 'post',
    body: JSON.stringify(graphQLParams),
    credentials: 'include',
  }).then(response => response.text()).then((responseBody) => {
    try {
      return JSON.parse(responseBody);
    } catch (error) {
      return responseBody;
    }
  });
}

ReactDOM.render(<GraphiQL fetcher={graphQLFetcher} style={{ height: '100vh' }} />, global.document.getElementById('graphiql'));
