import styles from './auth.module.scss';
import Header from './comp/header';
import { useSignupHooks } from './hooks/signup';
import { useUserContext } from './hooks/UserContext';

function Signup() {
    const { setRecheckAuth } = useUserContext();

    const {
        newUserName, 
        newUserPass,
        handleNewUserNameChange,
        handleNewUserPassChange,
        handleRegistration
    } = useSignupHooks(setRecheckAuth);

    return (
        <div className={styles.root}>
            <Header />
            <section className={styles.signupPage}>
                <div className={styles.infoSection}>
                    <div className={styles.infoContent}>
                    <h4>If the button turns green, your auth token is valid</h4>
                    </div>
                </div>
                <div className={styles.signupContainer}>
                    <form className={styles.signupForm} onSubmit={(e) => handleRegistration(e)}>
                    <h1>Sign Up</h1>
                    <p>Already have an account? <a href="/login">Login</a></p>
                        <input 
                            className={styles.inputField} 
                            type="text" 
                            placeholder="Username" 
                            value={newUserName} 
                            onChange={handleNewUserNameChange}
                        />
                        <input 
                            className={styles.inputField} 
                            type="password" 
                            placeholder="Password" 
                            value={newUserPass} 
                            onChange={handleNewUserPassChange}
                        />
                        <button type="submit" className={styles.signupButton}>Sign Up</button>
                    </form>
                </div>
            </section>
        </div>
    );
}

export default Signup;
