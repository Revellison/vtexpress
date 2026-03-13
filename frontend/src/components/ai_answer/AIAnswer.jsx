import styles from './AIAnswer.module.css';

const AIAnswer = ({ summary }) => {
    if (!summary) {
        return null;
    }

    return (
        <section className={styles.wrapper}>
            <h3>AI summary</h3>
            <p>{summary}</p>
        </section>
    );
};

export default AIAnswer;
