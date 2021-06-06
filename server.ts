import {config as envConfig} from 'dotenv'

const main = async () =>
{
    envConfig()
    console.log(process.env.MESSAGE)
}


main().catch(console.dir)