import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import LiveFeed from './components/LiveFeed';
import IncidentDetail from './components/IncidentDetail';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<LiveFeed />} />
          <Route path="/incidents/:id" element={<IncidentDetail />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
