import styles from './header.module.scss';
import { useState, useEffect } from 'react';
import GithubLogo from './../assets/github-mark-white.png';
import { ReactComponent as No } from './../assets/no.svg';
import { ReactComponent as Yes } from './../assets/yes.svg';
import { useUserContext } from '../hooks/UserContext';  
import { Link } from 'react-router-dom';

function Header() {
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const [isActive, setIsActive] = useState<boolean | null>(null);
    const { recheckAuth } = useUserContext();
    const [isLoggingOut, setIsLoggingOut] = useState(false);

    const handleLogout = async () => {
        try {
            const response = await fetch("http://localhost:8080/api/go/other/logout", {
                method: 'POST', 
                credentials: 'include'
            });
            if (response.status !== 200) {
                throw new Error("Logout failed");
            }
            setIsLoggingOut(false);
            setIsActive(false);
        } catch (error) {
            console.error("Logout failed:", error);
        }
    };
    useEffect(() => {
        if (isLoggingOut) {
            handleLogout();
        }
    }, [isLoggingOut]);

    const toggleDropdown = () => {
        setIsDropdownOpen(!isDropdownOpen);
    }
    const fetchPing = async () => {
        try {
            const response = await fetch("http://localhost:8080/api/go/other/valid_token", {
                credentials: 'include'
            });
            if (response.status !== 200) {
                throw new Error("Unauthorized");
            }
            const data = await response.json();
            console.log(data);
            return data;  
        } catch (error) {
            console.error("Failed to ping:", error);
            throw error;
        }
    };

    useEffect(() => {
        const checkActive = async () => {
            try {
                const data = await fetchPing();
                if (data.status === "success") { 
                    setIsActive(true);
                } else {
                    setIsActive(false);
                }
            } catch (error) {
                setIsActive(false); 
            }
        };
    
        checkActive();
    }, [recheckAuth]);

    return (
        <div className={styles.stickyHeader}>
            <div className={styles.headerContainer}>
                <Link to="/">
                    {isActive === null 
                        ? <h4>Loading...</h4> 
                        : <h4>{isActive ? <Yes title="mainLogo" className={styles.mainLogo} /> : <No title="mainLogo" className={styles.mainLogo} />}</h4>}
                </Link>
            </div>
            <div className={styles.headerContainer}>
            </div>   
            <div className={styles.headerContainer}>
                <div className={styles.authLinks}>
                    {isActive === null ? (
                        <h6>Loading...</h6>
                    ) : isActive ? (
                        <a href="#" onClick={() => setIsLoggingOut(true)}>
                            <h6>Logout</h6>
                        </a>
                    ) : (
                        <>
                            <Link to="/login">
                                <h6>Login</h6>
                            </Link>
                            <Link to="/signup">
                                <h6>Sign Up</h6>
                            </Link>
                        </>
                    )}
                </div>
                <a href="https://github.com//" target="_blank" rel="noopener noreferrer">
                    <img src={GithubLogo} className={styles.smallLogo} />
                </a>
            </div>
        </div>
    );
}

export default Header;
