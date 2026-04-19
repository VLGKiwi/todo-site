import { FC } from "react"
import { TodoItemProps } from "./TodoItem.types"
import styles from './TodoItem.module.scss'

export const TodoItem: FC<TodoItemProps> = ({ todo }) => {
	return (
		<div className={styles.container}>
			<h2 className={styles.title}>{todo.title}</h2>
			<p className={styles.description}>{todo.description}</p>
			<p className={styles.complete}>{todo.completed ? 'Выполнено' : 'В работе'}</p>
		</div>
	)
}
