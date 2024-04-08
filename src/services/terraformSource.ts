import axios from 'axios';

type Versions = Array<{
    version: string
    protocols: Array<string>
    platforms: Array<{
        os: string
        arch: string
    }>
}>

type Package = {
    download_url: string
}

const getVersions = async (namespace: string, type: string): Promise<Versions | undefined> => {
    try {
        const {data} = await axios.get(`https://registry.terraform.io/v1/providers/${namespace}/${type}/versions`);
        return data?.versions
    } catch (error) {
        if (axios.isAxiosError(error)) {
            console.log(error)
        } else {
            console.log(error)
        }
    }
}

const getPackage = async (namespace: string, type: string, version: string, os: string, arch: string): Promise<Package | undefined> => {
    try {
        const {data} = await axios.get(`https://registry.terraform.io/v1/providers/${namespace}/${type}/${version}/download/${os}/${arch}`);
        return data
    } catch (error) {
        if (axios.isAxiosError(error)) {
            console.log(error)
        } else {
            console.log(error)
        }
    }
}

export {
    getVersions,
    getPackage
}