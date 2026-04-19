import { FC } from "react"
import { ButtonProps } from "./Button.types"
import styles from './Button.module.scss'

export const Button: FC<ButtonProps> = ({
	children,
	...props
}) => {
	return (
		<button {...props} className={styles.button}>
			{children}
		</button>
	)
}
