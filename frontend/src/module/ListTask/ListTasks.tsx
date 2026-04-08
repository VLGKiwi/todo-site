'use client'

import { useTodos } from "@/hooks/useTodos";

export const ListTask = () => {

	const { data: todos } = useTodos();

	return (
		<div>
			{todos?.map((todo) => (
				<div key={todo.id}>
					<h1>{todo.title}</h1>
					<p>{todo.description}</p>
					<p>{todo.completed ? 'Completed' : 'Not Completed'}</p>
				</div>
			))}
		</div>
	)
}
