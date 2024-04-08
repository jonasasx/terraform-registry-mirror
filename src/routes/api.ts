import {Router} from 'express';
import jetValidator from 'jet-validator';
import {IReq, IRes} from "@src/routes/types/express/misc";
import HttpStatusCodes from "@src/constants/HttpStatusCodes";
import {getPackage, getVersions} from "@src/services/terraformSource";


// **** Variables **** //

const apiRouter = Router(),
    validate = jetValidator();

apiRouter.get(`/:hostname/:namespace/:type/index.json`, async (req: IReq, res: IRes) => {
    const versions = await getVersions(req.params['namespace'], req.params['type'])
    return res.status(HttpStatusCodes.OK).json({versions: versions?.map(v => v.version)});
})

apiRouter.get(`/:hostname/:namespace/:type/:version.json`, async (req: IReq, res: IRes) => {
    const versions = await getVersions(req.params['namespace'], req.params['type'])
    const version = versions?.find(v => v.version == req.params['version'])!!

    return res.status(HttpStatusCodes.OK).json({
        archives: await version.platforms.reduce<any>( // { [key: string]: string }
            async (accumulator, currentValue) => {
                return {
                    ...await accumulator,
                    [`${currentValue.os}_${currentValue.arch}`]: {
                        "url": (await getPackage(
                            req.params['namespace'],
                            req.params['type'],
                            req.params['version'],
                            currentValue.os,
                            currentValue.arch))?.download_url
                    }
                }
            },
            {},
        )
    });
})


// **** Export default **** //

export default apiRouter;
