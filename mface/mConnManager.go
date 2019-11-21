/**
* @Author: Cooper
* @Date: 2019/11/18 14:38
 */

package mface

type MConnManager interface {
	Add(string , MConnection) error
	Remove(string)
	CleanAll()
	Get(string) (MConnection , error)
	Len() uint32
}
