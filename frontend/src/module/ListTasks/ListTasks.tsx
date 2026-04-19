'use client'

import { TodoItem } from "@/components/TodoItem/TodoItem";
import { useTodos } from "@/hooks/useTodos";
import styles from './ListTasks.module.scss'

export const ListTasks = () => {

	const { data: todos } = useTodos();

	return (
		<div className={styles.container}>
			{todos?.map((todo) => (
				<TodoItem key={todo.id} todo={todo} />
			))}
		</div>
	)
}
