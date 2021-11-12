import express from 'express'
import winston from 'winston'

const { createLogger, format, transports } = winston

// TODO: Manage logger configuration from file or env variables
const logger = createLogger({
	format: format.combine(
		format.timestamp(),
		format.simple()
	),
	transports: [
		new transports.Console({
			format: format.combine(
				format.timestamp(),
				format.colorize(),
				format.simple()
			)
		}),
	]
})

// TODO: Add better support for configuration of port
const port = process.env.PORT || '8080'

const app = express()
app.get("/", (_, res) => {
	res.json({
		message: 'pippo'
	})
	res.end()
})

app.all("/api", (_req, res) => {
	// TODO: Entry point for all http esb flows
	res.send("done")
})

app.listen(port, () => {
	logger.log({
		level: 'info',
		message: `start server listening on port: ${port}`
	})
})
