/**
* @Author: Cooper
* @Date: 2019/11/21 19:28
 */

package mfaceV2

type MServer interface {
	Start(MConfig)
	Stop()
}