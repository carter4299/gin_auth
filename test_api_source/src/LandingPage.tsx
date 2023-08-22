import styles from './app.module.scss';
import Header from './comp/header';

function LandingPage() {
    return (
        <div className={styles.root}>
            <Header />
            <section className={styles.container}>
            <h4>Use Login / Signup / Logout in Header</h4>
            </section>
        </div>
    );
}

export default LandingPage;


