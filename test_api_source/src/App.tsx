
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LandingPage from './LandingPage';
import Signup from './Signup';
import Login from './Login';
import { useState } from 'react';
import { UserContext } from './hooks/UserContext';

function App() {
    const [recheckAuth, setRecheckAuth] = useState(false);
    return (
        <UserContext.Provider value={{ recheckAuth, setRecheckAuth }}>
            <Router>
                    <Routes>
                        <Route path="/" element={<LandingPage/>} />
                        <Route path="/signup" element={<Signup/>} />
                        <Route path="/login" element={<Login/>} />
                    </Routes>
            </Router>
        </UserContext.Provider>
    );
}

export default App;


