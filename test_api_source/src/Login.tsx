import styles from './auth.module.scss';
import Header from './comp/header';
import { useLoginHooks } from './hooks/login';
import { useUserContext } from './hooks/UserContext';

function Login() {
    const { setRecheckAuth } = useUserContext();

    const {
        userName,
        userPass,
        handleUserNameChange,
        handleUserPassChange,
        handleSubmit
    } = useLoginHooks(setRecheckAuth);

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
                    <form className={styles.signupForm} onSubmit={(e) => handleSubmit(e)}>
                    <h1>Sign in to your Account</h1>
                        <input 
                            className={styles.inputField} 
                            type="text" 
                            placeholder="Username" 
                            value={userName} 
                            onChange={handleUserNameChange}
                        />
                        <input 
                            className={styles.inputField} 
                            type="password" 
                            placeholder="Password" 
                            value={userPass} 
                            onChange={handleUserPassChange}
                        />
                        <button type="submit" className={styles.signupButton}>Sign In</button>
                    </form>
                </div>
            </section>
        </div>
    );
}

export default Login;
